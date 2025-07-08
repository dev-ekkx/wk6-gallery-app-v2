import {Component, inject, signal} from '@angular/core';
import {ImageService} from '../../services/image/image';
import {take} from 'rxjs';

@Component({
  selector: 'app-image-upload',
  imports: [],
  templateUrl: './image-upload.html',
  styleUrl: './image-upload.css'
})
export class ImageUpload {
protected imageService = inject(ImageService);
  protected images = signal<{ file: File, url: string }[]>([]);
  protected isUploading = signal(false);

  protected onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      this.processFiles(input.files);
    }
  }

  protected onDragOver(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
    const dropZone = event.currentTarget as HTMLElement;
    dropZone.classList.add('dragover');
  }

  protected onDragLeave(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
    const dropZone = event.currentTarget as HTMLElement;
    dropZone.classList.remove('dragover');
  }

  protected onDrop(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
    const dropZone = event.currentTarget as HTMLElement;
    dropZone.classList.remove('dragover');

    if (event.dataTransfer?.files && event.dataTransfer.files.length > 0) {
      this.processFiles(event.dataTransfer.files);
    }
  }

  protected removeImage(index: number): void {
    const newImagesList = this.images().filter((_, i) => i !== index);
    this.images.set(newImagesList);
  }

  private processFiles(files: FileList): void {
    Array.from(files).forEach(file => {
      if (file.type.startsWith('image/')) {
        const reader = new FileReader();
        reader.onload = (e: ProgressEvent<FileReader>) => {
          if (e.target?.result) {
            const image = {
              file,
              url: e.target.result as string
            }
            this.images.set([...this.images(), image])
          }
        };
        reader.readAsDataURL(file);

      } else {
        alert("Only image files are allowed.");
      }
    });
  }

  protected uploadToS3(): void {
    this.isUploading.set(true);
    const imagesArray = this.images().map(image => image.file);
    console.log(imagesArray);
this.imageService.uploadImages(imagesArray).pipe(take(1)).subscribe({
  next: (result) => {
    console.log('Upload successful:', result);
    this.images.set([]);
  },
  error: (error) => {
    console.error('Upload failed:', error);
    alert('Failed to upload images. Please try again.');
  },
  complete: () => {
    this.isUploading.set(false);
  }
})

  }
}
