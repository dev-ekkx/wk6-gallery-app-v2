import {Component, signal} from '@angular/core';
import { RouterOutlet } from '@angular/router';
import {ImageUpload} from './components/image-upload/image-upload';
import {ImageList} from './components/image-list/image-list';

@Component({
  selector: 'app-root',
  imports: [ImageUpload, ImageList],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  protected title = 'frontend';



}
