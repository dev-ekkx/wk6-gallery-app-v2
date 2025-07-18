import {inject, Injectable, signal} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Image, ImageResponse, ImageUploadInterface} from '../../interfaces/interfaces';

@Injectable({
  providedIn: 'root'
})
export class ImageService {
  protected http = inject(HttpClient);
  public   images = signal<Image[]>([]);

public uploadImages(images: ImageUploadInterface[]) {
  const formData = new FormData();
 images.forEach((img, i) => {
    formData.append(`images[${i}].file`, img.file);
    formData.append(`images[${i}].description`, img.description);
  });
  return this.http.post<{ message: string }>('/api/upload', formData);
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
