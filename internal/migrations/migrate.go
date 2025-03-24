package migrations

import (
	"log"

	"invservice/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		return err
	}

	enumQuery := `
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'resource_status') THEN
        CREATE TYPE resource_status AS ENUM ('available', 'booked', 'unavailable');
    END IF;
END$$;
`
	if err := db.Exec(enumQuery).Error; err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Resource{}); err != nil {
		return err
	}

	if err := db.Exec(`DROP TRIGGER IF EXISTS update_resources_updated_at ON resources;`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;
`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
CREATE TRIGGER update_resources_updated_at
BEFORE UPDATE ON resources
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
`).Error; err != nil {
		return err
	}

	log.Println("Migration completed successfully!")
	return nil
}
