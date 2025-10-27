package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx stdlib driver for database/sql
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// NOTE: no comma after password
		dsn = "postgres://postgres:wildan123@127.0.0.1:5432/topup?sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	if err := seedAll(ctx, db); err != nil {
		log.Fatalf("seeding failed: %v", err)
	}

	fmt.Println("✅ Seeds applied successfully.")
}

func seedAll(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	now := time.Now()

	// 1) user_levels
	if err := exec(tx, `
INSERT INTO user_levels (name, description, created_at, updated_at)
SELECT $1::text, $2::text, $3, $3
WHERE NOT EXISTS (SELECT 1 FROM user_levels WHERE name=$1::text);
`, "Basic", "Default user level", now); err != nil {
		return err
	}
	if err := exec(tx, `
INSERT INTO user_levels (name, description, created_at, updated_at)
SELECT $1::text, $2::text, $3, $3
WHERE NOT EXISTS (SELECT 1 FROM user_levels WHERE name=$1::text);
`, "Silver", "Mid-tier user level", now); err != nil {
		return err
	}
	if err := exec(tx, `
INSERT INTO user_levels (name, description, created_at, updated_at)
SELECT $1::text, $2::text, $3, $3
WHERE NOT EXISTS (SELECT 1 FROM user_levels WHERE name=$1::text);
`, "Gold", "Premium user level", now); err != nil {
		return err
	}

	// 2) admin user (uses password_hash and user_level_id FK → Basic)
	if err := exec(tx, `
INSERT INTO users (name, email, password_hash, user_level_id, balance, role, whatsapp, created_at, updated_at)
SELECT $1::text, $2::text, $3::text,
       (SELECT id FROM user_levels WHERE name='Basic'),
       $4::numeric, $5::text, $6::text, $7, $7
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email=$2::text);
`, "Admin", "admin@topup.local",
		// bcrypt("admin123") — change later
		"$2a$10$V4xQz0KX4m0p0Lw6Qx6u8uI1k0mF0cZkqT0w1k3UoUSe3s6Uq9C9y",
		0, "admin", "081234567890", now); err != nil {
		return err
	}

	// 3) menus
	for _, m := range []string{"Home", "Pulsa & Data", "Games"} {
		if err := exec(tx, `
INSERT INTO menus (name, created_at, updated_at)
SELECT $1::text, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM menus WHERE name=$1::text);
`, m); err != nil {
			return err
		}
	}

	// 4) providers (uses ref, not code)
	type Prov struct{ Name, Ref string }
	for _, p := range []Prov{
		{"Digiflazz", "digiflazz"},
		{"XPay", "xpay"},
	} {
		if err := exec(tx, `
INSERT INTO providers (name, ref, created_at, updated_at)
SELECT $1::text, $2::text, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM providers WHERE ref=$2::text);
`, p.Name, p.Ref); err != nil {
			return err
		}
	}

	// 5) categories (per your entity fields)
	type Cat struct {
		Name, Type, MenuName, ProviderRef, Slug, Status, Description, InputType, ImgUrl string
		IsLogin                                                                         bool
	}
	cats := []Cat{
		{"Pulsa Prabayar", "prabayar", "Pulsa & Data", "digiflazz", "pulsa-prabayar", "active", "Pulsa prabayar", "phone", "https://example.com/img/cat-pulsa.png", false},
		{"Data Internet", "prabayar", "Pulsa & Data", "digiflazz", "data-internet", "active", "Paket data", "phone", "https://example.com/img/cat-data.png", false},
		{"PDAM", "pascabayar", "Home", "xpay", "pdam", "active", "Pembayaran PDAM", "customer_id", "https://example.com/img/cat-pdam.png", true},
		{"Online Games", "prabayar", "Games", "digiflazz", "online-games", "active", "Top-up game online", "game_id", "https://example.com/img/cat-games.png", false},
	}
	for _, c := range cats {
		if err := exec(tx, `
INSERT INTO categories (name, type, menu_id, provider_id, slug, status, description, input_type, img_url, is_login, created_at, updated_at)
SELECT $1::text, $2::text,
       (SELECT id FROM menus WHERE name=$3::text),
       (SELECT id FROM providers WHERE ref=$4::text),
       $5::text, $6::text, $7::text, $8::text, $9::text, $10::boolean, $11, $11
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE slug=$5::text);
`, c.Name, c.Type, c.MenuName, c.ProviderRef, c.Slug, c.Status, c.Description, c.InputType, c.ImgUrl, c.IsLogin, now); err != nil {
			return fmt.Errorf("insert category %q: %w", c.Name, err)
		}
	}

	// 6) products (updated schema: sku_code, seller_name, stock, img_url, start_off, end_off)
	type Prod struct {
		Name, Sku, Seller, CatSlug, ProviderRef, Status, Desc, Img, StartOff, EndOff string
		Stock                                                                        int64
		BasePrice                                                                    float64
	}
	prods := []Prod{
		// Pulsa / Data / PDAM
		{"Pulsa 25K", "PULSA-25K", "TopUp Demo", "pulsa-prabayar", "digiflazz", "active", "Pulsa prabayar nominal 25.000", "https://example.com/img/pulsa25k.png", "", "", 99999, 25000},
		{"Data 5GB", "DATA-5GB-30D", "TopUp Demo", "data-internet", "digiflazz", "active", "Paket data 5GB 30 hari", "https://example.com/img/data5gb.png", "", "", 99999, 30000},
		{"PDAM Kota A", "PDAM-KOTA-A", "TopUp Demo", "pdam", "xpay", "active", "Pembayaran tagihan PDAM Kota A", "https://example.com/img/pdam.png", "", "", 99999, 50000},

		// Online Games
		{"Mobile Legends Diamonds 86", "ML-86", "TopUp Demo", "online-games", "digiflazz", "active", "Top-up Mobile Legends 86 Diamonds", "https://example.com/img/ml-86.png", "00:00", "23:59", 99999, 20000},
		{"Free Fire Diamonds 100", "FF-100", "TopUp Demo", "online-games", "digiflazz", "active", "Top-up Free Fire 100 Diamonds", "https://example.com/img/ff-100.png", "00:00", "23:59", 99999, 25000},
		{"PUBG Mobile UC 60", "PUBG-60", "TopUp Demo", "online-games", "xpay", "active", "Top-up PUBG Mobile 60 UC", "https://example.com/img/pubg-60.png", "00:00", "23:59", 99999, 30000},
	}

	for _, p := range prods {
		if err := exec(tx, `
INSERT INTO products
(name, sku_code, seller_name, category_id, provider_id, status, stock, base_price, description, img_url, start_off, end_off, created_at, updated_at)
SELECT
  $1::text,  -- name
  $2::text,  -- sku_code
  $3::text,  -- seller_name
  (SELECT id FROM categories WHERE slug=$4::text), -- category_id
  (SELECT id FROM providers  WHERE ref=$5::text),  -- provider_id
  $6::text,               -- status
  $7::bigint,             -- stock
  $8::numeric,             -- base price
  $9::text,               -- description
  $10::text,               -- img_url
  $11::text,              -- start_off
  $12::text,              -- end_off
  $13, $13                -- created_at, updated_at
WHERE NOT EXISTS (SELECT 1 FROM products WHERE sku_code=$2::text);
`, p.Name, p.Sku, p.Seller, p.CatSlug, p.ProviderRef, p.Status, p.Stock, p.BasePrice, p.Desc, p.Img, p.StartOff, p.EndOff, now); err != nil {
			return fmt.Errorf("insert product %q: %w", p.Name, err)
		}
	}

	// 7) prices (amount column is numeric)
	type PriceRow struct {
		ProductName, Level string
		Amount             float64
	}
	priceRows := []PriceRow{
		{"Pulsa 25K", "Basic", 25000}, {"Pulsa 25K", "Silver", 24500}, {"Pulsa 25K", "Gold", 24000},
		{"Data 5GB", "Basic", 50000}, {"Data 5GB", "Silver", 49000}, {"Data 5GB", "Gold", 48000},
		{"PDAM Kota A", "Basic", 3500}, {"PDAM Kota A", "Silver", 3300}, {"PDAM Kota A", "Gold", 3100},

		{"Mobile Legends Diamonds 86", "Basic", 20000},
		{"Mobile Legends Diamonds 86", "Silver", 19500},
		{"Mobile Legends Diamonds 86", "Gold", 19000},
		{"Free Fire Diamonds 100", "Basic", 15000},
		{"Free Fire Diamonds 100", "Silver", 14750},
		{"Free Fire Diamonds 100", "Gold", 14500},
		{"PUBG Mobile UC 60", "Basic", 12000},
		{"PUBG Mobile UC 60", "Silver", 11800},
		{"PUBG Mobile UC 60", "Gold", 11600},
	}
	for _, r := range priceRows {
		if err := exec(tx, `
INSERT INTO prices (product_id, user_level_id, amount, created_at, updated_at)
SELECT p.id, l.id, $1::numeric, $2, $2
FROM products p, user_levels l
WHERE p.name=$3::text AND l.name=$4::text
AND NOT EXISTS (SELECT 1 FROM prices WHERE product_id=p.id AND user_level_id=l.id);
`, r.Amount, now, r.ProductName, r.Level); err != nil {
			return fmt.Errorf("insert price %q/%q: %w", r.ProductName, r.Level, err)
		}
	}

	// 8) payment_methods (type/name/img_url/provider_id/fee/percent)
	if err := exec(tx, `
INSERT INTO payment_methods (type, name, img_url, provider_id, fee, percent, created_at, updated_at)
SELECT $1::text, $2::text, $3::text,
       (SELECT id FROM providers WHERE ref=$4::text),
       $5::numeric, $6::numeric, $7, $7
WHERE NOT EXISTS (SELECT 1 FROM payment_methods WHERE name=$2::text);
`, "bank", "Bank Transfer (BCA)", "https://example.com/img/bca.png", "digiflazz", 1000.0, 0.0, now); err != nil {
		return err
	}
	if err := exec(tx, `
INSERT INTO payment_methods (type, name, img_url, provider_id, fee, percent, created_at, updated_at)
SELECT $1::text, $2::text, $3::text,
       (SELECT id FROM providers WHERE ref=$4::text),
       $5::numeric, $6::numeric, $7, $7
WHERE NOT EXISTS (SELECT 1 FROM payment_methods WHERE name=$2::text);
`, "ewallet", "OVO", "https://example.com/img/ovo.png", "digiflazz", 0.0, 1.0, now); err != nil {
		return err
	}

	// 9) banners
	if err := exec(tx, `
INSERT INTO banners (img_url, created_at, updated_at)
SELECT $1::text, $2, $2
WHERE NOT EXISTS (SELECT 1 FROM banners WHERE img_url=$1::text);
`, "https://example.com/b/sept.png", now); err != nil {
		return err
	}

	// 10) settings
	settings := [][2]string{
		{"site_name", "TopUp Demo"},
		{"site_logo", "https://example.com/logo.png"},
		{"support_email", "support@topup.local"},
	}
	for _, s := range settings {
		if err := exec(tx, `
INSERT INTO settings (name, value, created_at, updated_at)
SELECT $1::text, $2::text, $3, $3
WHERE NOT EXISTS (SELECT 1 FROM settings WHERE name=$1::text);
`, s[0], s[1], now); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func exec(tx *sql.Tx, q string, args ...any) error {
	_, err := tx.Exec(q, args...)
	return err
}
