export function formatDate(value: string | null | undefined) {
  if (!value) return 'N/A';
  return new Intl.DateTimeFormat(undefined, {
    dateStyle: 'medium',
    timeStyle: value.includes('T') ? 'short' : undefined,
  }).format(new Date(value));
}

export function relativeTime(value: string) {
  const diff = Date.now() - new Date(value).getTime();
  const minutes = Math.round(diff / 60000);
  if (minutes < 1) return 'just now';
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.round(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  return `${Math.round(hours / 24)}d ago`;
}

export function titleCase(input: string) {
  return input.replace(/_/g, ' ').replace(/\b\w/g, (char) => char.toUpperCase());
}
