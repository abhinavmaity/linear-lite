DELETE FROM labels
WHERE (name, color) IN (
  ('bug', '#EF4444'),
  ('feature', '#22C55E'),
  ('design', '#EAB308'),
  ('infra', '#2563EB'),
  ('frontend', '#EC4899')
);
