import {inject, Injectable, signal} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Image, ImageResponse, ImageUploadInterface} from '../../interfaces/interfaces';

@Injectable({
  providedIn: 'root'
})
export class ImageService {
  protected http = inject(HttpClient);
  public   images = signal<Image[]>([]);

  // host = 'http://localhost:8080/api';
  host = "/api";

public uploadImages(images: ImageUploadInterface[]) {
  const formData = this.buildFormData(images);

    this.logFormData(formData);

  return this.http.post<{ message: string }>(
    `${this.host}/upload`,
    formData
  );
}

 public getImages(startAfter?: string) {
    let url = `${this.host}/images`;
    if (startAfter) {
      url += `?startAfter=${encodeURIComponent(startAfter)}`;
    }
    return this.http.get<ImageResponse>(url);
  }

  public deleteImage(key: string) {
    return this.http.delete<{ message: string }>(`${this.host}/api/images/${encodeURIComponent(key)}`);
  }
  
  
  private buildFormData(images: ImageUploadInterface[]): FormData {
    const formData = new FormData();
    images.forEach((img, index) => {
      formData.append(`images[${index}].file`, img.file);
      formData.append(`images[${index}].description`, img.description);
    });
    return formData;
  }

  private logFormData(formData: FormData): void {
  for (const [key, value] of formData.entries()) {
    console.log(`${key}:`, value);
  }
}

}

