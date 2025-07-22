import {Component, computed, inject, signal} from '@angular/core';
import {ImageService} from '../../services/image/image';
import {take} from 'rxjs';
import { ImageUploadInterface } from '../../interfaces/interfaces';

@Component({
  selector: 'app-image-upload',
  imports: [],
  templateUrl: './image-upload.html',
})
export class ImageUpload {
protected imageService = inject(ImageService);
  protected images = signal<ImageUploadInterface[]>([]);
  protected isUploading = signal(false);
  protected isDragOver = signal(false);

  checkDescriptionsValidity = computed(() => {
    return this.images().every(img => img.description.trim().length > 2);
  });
  
   onDragOver(event: DragEvent) {
     this.isDragOver.set(true);
     event.preventDefault();
     console.log(this.isDragOver())
  }

  onDragLeave(event: DragEvent) {
    this.isDragOver.set(false); 
    event.preventDefault();
   }

   onDrop(event: DragEvent) {
    console.log(event)
    this.isDragOver.set(false);
    event.preventDefault();
    if (event.dataTransfer?.files) {
      this.processFiles(event.dataTransfer.files);
    }
  }

  onInputChange(event: Event, index: number) {
 const value = (event.target as HTMLInputElement).value;

  this.images.update((images) =>
    images.map((img, i) =>
      i === index ? { ...img, description: value } : img
    ))
  }


 onFileChange(event: Event) {
    const input = event.target as HTMLInputElement;
    if (input.files) {
      this.processFiles(input.files);
    }
  }

    processFiles(fileList: FileList) {
    Array.from(fileList).forEach((file) => {

      if (!file.type.startsWith('image/')) {
        alert(`${file.name} is not an image`);
        return;
      }
      const reader = new FileReader();
      reader.onload = () => {
        this.images.update(images => [
          ...images,
          { file, url: reader.result as string, description: '' }
        ]);
      };
      reader.readAsDataURL(file);
    });
  }

  onRemove(index: number) {
    this.images.update(images => images.filter((_, i) => i !== index));
  }



  upload() {
    this.isUploading.set(true);
    this.imageService.uploadImages(this.images())
      .pipe(take(1))
      .subscribe({
        next: (response) => {
          this.isUploading.set(false);
          console.log(response.message);
          this.images.set([]);
        }
        , error: (error) => {
          this.isUploading.set(false);
          alert('Upload failed: ' + error.message);
        }
        , complete: () => {
          this.isUploading.set(false);
          this.imageService.getImages().pipe(take(1)).subscribe({
            next: (response) => {
              this.imageService.images.set(response.images);
            },
            error: (error) => {
              alert('Failed to fetch images: ' + error.message);
            }
          });
        }
      })
  }
}
