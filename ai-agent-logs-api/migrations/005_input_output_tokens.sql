-- Replace token_usage text with explicit input_token / output_token integers.

ALTER TABLE session_prompts
    ADD COLUMN IF NOT EXISTS input_token INT NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS output_token INT NOT NULL DEFAULT 0;

DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'session_prompts' AND column_name = 'token_usage'
  ) THEN
    UPDATE session_prompts sp
    SET
      input_token = CASE
        WHEN sp.token_usage ~* 'in:\s*\d+' THEN
          COALESCE((regexp_match(sp.token_usage, 'in:\s*(\d+)', 'i'))[1]::int, 0)
        ELSE 0
      END,
      output_token = CASE
        WHEN sp.token_usage ~* 'out:\s*\d+' THEN
          COALESCE((regexp_match(sp.token_usage, 'out:\s*(\d+)', 'i'))[1]::int, 0)
        WHEN sp.token_usage ~ '^\s*\d+\s*$' THEN btrim(sp.token_usage)::int
        WHEN sp.token_usage ~* 'in:\s*\d+' THEN 0
        ELSE COALESCE(parse_token_usage(sp.token_usage), 0)
      END
    WHERE btrim(COALESCE(sp.token_usage, '')) <> ''
      AND lower(btrim(sp.token_usage)) <> 'unknown';

    ALTER TABLE session_prompts DROP COLUMN token_usage;
  END IF;
END $$;

DROP FUNCTION IF EXISTS parse_token_usage(TEXT);
