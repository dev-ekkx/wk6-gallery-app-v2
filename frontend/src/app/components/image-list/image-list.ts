import { Component } from '@angular/core';

@Component({
  selector: 'app-image-list',
  imports: [],
  templateUrl: './image-list.html',
  styleUrl: './image-list.css'
})
export class ImageList {
  images: {url: string, file: File}[] = []

  protected formatFileSize(size: number): string {
    if (size < 1024) return `${size} B`;
    if (size < 1048576) return `${(size / 1024).toFixed(2)} KB`;
    return `${(size / 1048576).toFixed(2)} MB`;
  }

  protected deleteImage(index: number): void {
    console.log(index)
  }
}
