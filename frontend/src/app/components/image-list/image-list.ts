import {Component, inject, OnDestroy, OnInit, signal} from '@angular/core';
import {ImageService} from '../../services/image/image';
import {Subject, takeUntil} from 'rxjs';
import {Image} from '../../interfaces/interfaces';

@Component({
  selector: 'app-image-list',
  imports: [],
  templateUrl: './image-list.html',
  styleUrl: './image-list.css'
})
export class ImageList implements OnInit, OnDestroy {
  protected imageService = inject(ImageService);

  images = signal<Image[]>([]);
  nextToken  = signal<string | null>(null);
  hasMore = signal(false);
  loading = signal(false);

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

  protected deleteImage(index: number): void {
    console.log(index)
  }

  ngOnDestroy() {
    this.destroy$.next()
    this.destroy$.complete()
  }
}
