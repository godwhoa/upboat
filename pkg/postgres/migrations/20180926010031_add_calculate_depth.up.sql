CREATE OR REPLACE FUNCTION calculate_depth(parent_id integer) 
RETURNS integer AS $$
DECLARE parent_depth INTEGER
BEGIN
        IF parent_id IS NULL THEN
            RETURN 0;
        END IF;
        SELECT depth INTO parent_depth FROM comments WHERE id = parent_id;
        RETURN parent_depth + 1;
END;
$$ LANGUAGE plpgsql;