-- App field for agent-turn logging + helper to parse token_usage strings.

ALTER TABLE sessions
    ADD COLUMN IF NOT EXISTS app TEXT NOT NULL DEFAULT 'cursor';

ALTER TABLE session_prompts
    ADD COLUMN IF NOT EXISTS app TEXT NOT NULL DEFAULT 'cursor';

CREATE OR REPLACE FUNCTION parse_token_usage(t TEXT) RETURNS INT AS $$
DECLARE
  total INT := 0;
  in_m TEXT[];
  out_m TEXT[];
  m TEXT[];
BEGIN
  IF t IS NULL OR btrim(t) = '' OR lower(btrim(t)) = 'unknown' THEN
    RETURN 0;
  END IF;

  IF t ~* 'in:\s*\d+' AND t ~* 'out:\s*\d+' THEN
    in_m := regexp_match(t, 'in:\s*(\d+)', 'i');
    out_m := regexp_match(t, 'out:\s*(\d+)', 'i');
    RETURN COALESCE(in_m[1]::int, 0) + COALESCE(out_m[1]::int, 0);
  END IF;

  IF t ~ '^\s*\d+\s*$' THEN
    RETURN btrim(t)::int;
  END IF;

  FOR m IN SELECT regexp_matches(t, '(\d+)', 'g') LOOP
    total := total + m[1]::int;
  END LOOP;
  RETURN total;
END;
$$ LANGUAGE plpgsql IMMUTABLE;
