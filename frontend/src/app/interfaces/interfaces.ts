export interface Image {
  key: string;
  size: number;
  etag: string;
  url: string;
  description: string;
}
export interface ImageResponse {
  images: Image[];
  nextToken?: string;
  isTruncated?: boolean;
}

export interface ImageUploadInterface {
  file: File;
  url: string;
  description: string;
}
