import { ApiError } from 'types/api';

export interface ParsedUiError {
  message: string;
  summary: string;
  fields: Record<string, string> | null;
}

const FIELD_LABELS: Record<string, string> = {
  name: 'Name',
  email: 'Email',
  password: 'Password',
  title: 'Title',
  description: 'Description',
  key: 'Project key',
  color: 'Color',
  project_id: 'Project',
  sprint_id: 'Sprint',
  assignee_id: 'Assignee',
  start_date: 'Start date',
  end_date: 'End date',
  status: 'Status',
  priority: 'Priority',
  label_ids: 'Labels',
};

function toSentenceCase(value: string) {
  if (!value) return value;
  return value.charAt(0).toUpperCase() + value.slice(1);
}

function fieldLabel(field: string) {
  return FIELD_LABELS[field] ?? field.replace(/_/g, ' ');
}

function normalizeFieldMessage(field: string, message: string) {
  const trimmed = (message || '').trim();
  if (!trimmed) {
    return `${fieldLabel(field)} is invalid.`;
  }

  if (trimmed.toLowerCase() === 'is required' || trimmed.toLowerCase() === 'this field is required.') {
    return `${fieldLabel(field)} is required.`;
  }

  if (trimmed.includes('must be less than or equal to')) {
    const match = trimmed.match(/(\d+)/);
    if (match) {
      return `${fieldLabel(field)} must be ${match[1]} characters or fewer.`;
    }
  }

  if (trimmed.includes('must match ^[A-Z0-9]{2,10}$')) {
    return 'Project key must be 2-10 uppercase letters or numbers.';
  }

  if (trimmed.includes('must match #RRGGBB')) {
    return 'Color must be a 6-digit hex code like #3B82F6.';
  }

  if (trimmed.includes('must be between 8 and 72 characters')) {
    return 'Password must be 8-72 characters long.';
  }

  if (trimmed.toLowerCase().includes('must be a valid email address')) {
    return 'Enter a valid email address.';
  }

  if (trimmed.toLowerCase() === 'already in use') {
    return `${fieldLabel(field)} is already in use.`;
  }

  if (trimmed === 'must be greater than or equal to start_date') {
    return 'End date must be on or after the start date.';
  }

  if (trimmed.includes('must use YYYY-MM-DD format')) {
    return `${fieldLabel(field)} must use YYYY-MM-DD format.`;
  }

  if (trimmed.includes('must be a valid UUID') || trimmed.includes('must be a valid ID')) {
    return `Select a valid ${fieldLabel(field).toLowerCase()}.`;
  }

  if (trimmed.includes('must not contain duplicates') || trimmed.includes('Duplicate values are not allowed')) {
    return `Remove duplicate ${fieldLabel(field).toLowerCase()}.`;
  }

  return toSentenceCase(trimmed.endsWith('.') ? trimmed : `${trimmed}.`);
}

function normalizeGeneralMessage(message: string, fallback: string) {
  const trimmed = (message || '').trim();
  if (!trimmed) return fallback;

  if (trimmed.toLowerCase() === 'invalid email or password') {
    return 'Email or password is incorrect.';
  }

  if (trimmed.toLowerCase() === 'authentication is required') {
    return 'Please sign in to continue.';
  }

  return toSentenceCase(trimmed.endsWith('.') ? trimmed : `${trimmed}.`);
}

export function parseUiError(error: unknown, fallback: string): ParsedUiError {
  if (!(error instanceof ApiError)) {
    return {
      message: fallback,
      summary: fallback,
      fields: null,
    };
  }

  const normalizedFields = error.fields
    ? Object.fromEntries(Object.entries(error.fields).map(([field, message]) => [field, normalizeFieldMessage(field, message)]))
    : null;

  const hasFieldErrors = Boolean(normalizedFields && Object.keys(normalizedFields).length > 0);
  const normalizedMessage = normalizeGeneralMessage(error.message, fallback);
  const summary = hasFieldErrors && error.code === 'validation_error'
    ? 'Please correct the highlighted fields and try again.'
    : normalizedMessage;

  return {
    message: normalizedMessage,
    summary,
    fields: normalizedFields,
  };
}

export function getBannerErrorMessage(error: unknown, fallback = 'Something went wrong. Please try again.') {
  return parseUiError(error, fallback).message;
}
