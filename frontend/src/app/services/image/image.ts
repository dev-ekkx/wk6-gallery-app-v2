import {inject, Injectable, signal} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Image, ImageResponse} from '../../interfaces/interfaces';

@Injectable({
  providedIn: 'root'
})
export class ImageService {
  protected http = inject(HttpClient);
  public   images = signal<Image[]>([]);

public uploadImages(files: File[]) {
  const form = new FormData();
  for (const file of files) {
    form.append('images', file);
  }
  return this.http.post<{ message: string }>('/api/upload', form);
}

 public getImages(startAfter?: string) {
    let url = '/api/images';
    if (startAfter) {
      url += `?startAfter=${encodeURIComponent(startAfter)}`;
    }
    return this.http.get<ImageResponse>(url);
  }

  public deleteImage(key: string) {
    return this.http.delete<{ message: string }>(`/api/images/${encodeURIComponent(key)}`);
  }
}
