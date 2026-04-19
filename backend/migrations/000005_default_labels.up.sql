DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM labels LIMIT 1) THEN
    INSERT INTO labels (name, color, description)
    VALUES
      ('bug', '#EF4444', 'Defects and regressions'),
      ('feature', '#22C55E', 'Net-new functionality'),
      ('design', '#EAB308', 'UI and UX improvements'),
      ('infra', '#2563EB', 'Infrastructure and DevOps'),
      ('frontend', '#EC4899', 'Frontend ownership');
  END IF;
END $$;
