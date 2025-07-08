import {inject, Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class ImageService {
  protected http = inject(HttpClient);

public uploadImages(files: File[]) {
  const form = new FormData();
  for (const file of files) {
    form.append('images', file); // key name must match backend
  }
  return this.http.post<{ message: string }>('/api/upload', form);
}
}
