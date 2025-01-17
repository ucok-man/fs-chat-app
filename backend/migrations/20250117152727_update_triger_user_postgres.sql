-- plpgsql-language-server:disable validation

-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION record_user_update()
RETURNS TRIGGER
AS
$$
BEGIN
    NEW.updated_at := now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER record_user_update_trigger
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION record_user_update();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS record_user_update_trigger ON users;
DROP FUNCTION IF EXISTS record_user_update();
-- +goose StatementEnd
