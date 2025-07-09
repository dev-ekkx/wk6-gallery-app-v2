export interface Image {
  key: string;
  size: number;
  etag: string;
  url: string;
}
export interface ImageResponse {
  images: Image[];
  nextToken?: string;
  isTruncated?: boolean;
}
