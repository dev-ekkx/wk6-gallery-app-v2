import {Component, effect, inject, OnDestroy, OnInit, signal} from '@angular/core';
import {ImageService} from '../../services/image/image';
import {Subject, take, takeUntil} from 'rxjs';
import {Image} from '../../interfaces/interfaces';

@Component({
  selector: 'app-image-list',
  imports: [],
  templateUrl: './image-list.html',
  styleUrl: './image-list.css'
})
export class ImageList implements OnInit, OnDestroy {
  protected imageService = inject(ImageService);
  images = signal<Image[]>(this.imageService.images());
  nextToken  = signal<string | null>(null);
  hasMore = signal(false);
  loading = signal(false);

  constructor() {
    effect(() => {
      if (this.imageService.images()) {
        this.images.set(this.imageService.images());
      }
    });  }

  protected destroy$ = new Subject<void>();

  ngOnInit(): void {
    this.loadImages();
  }

  loadImages(): void {
    this.loading.set(true);
    this.imageService.getImages(this.nextToken() ?? undefined).pipe(takeUntil(this.destroy$)).subscribe(
      {
        next: res => {
          this.images.set(res.images);
          this.nextToken.set(res.nextToken ?? null);
          this.hasMore.set(!!res.isTruncated);
        },
        error: () => {
          alert('Failed to load images');
        },
        complete: () => {
          this.loading.set(false);
        }
      }
    )
  }

  loadMore(): void {
    if (this.nextToken() && !this.loading()) {
      this.loadImages();
    }
  }

  protected formatFileSize(size: number): string {
    if (size < 1024) return `${size} B`;
    if (size < 1048576) return `${(size / 1024).toFixed(2)} KB`;
    return `${(size / 1048576).toFixed(2)} MB`;
  }

  protected deleteImage(key: string) {
    if (!confirm('Are you sure you want to delete this image?')) return;

    this.imageService.deleteImage(key).pipe(take(2)).subscribe({
      next: () => {
        const newImages = this.images().filter(img => img.key !== key);
        this.imageService.images.set(newImages);
      },
      error: err => {
        console.error('Failed to delete image', err);
      }
    });
  }


  ngOnDestroy() {
    this.destroy$.next()
    this.destroy$.complete()
  }
}
