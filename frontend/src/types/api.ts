export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface CollectionResponse<T> {
  items: T[];
  pagination: PaginationMeta;
}

export interface SingleResponse<T> {
  data: T;
}

export interface ErrorEnvelope {
  error: {
    code: string;
    message: string;
    fields?: Record<string, string>;
    request_id?: string;
  };
}

export class ApiError extends Error {
  code: string;
  status: number;
  fields?: Record<string, string>;
  requestId?: string;

  constructor(status: number, envelope: ErrorEnvelope['error']) {
    super(envelope.message);
    this.name = 'ApiError';
    this.status = status;
    this.code = envelope.code;
    this.fields = envelope.fields;
    this.requestId = envelope.request_id;
  }
}
