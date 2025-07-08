import {Component, signal} from '@angular/core';
import { RouterOutlet } from '@angular/router';
import {ImageUpload} from './components/image-upload/image-upload';

@Component({
  selector: 'app-root',
  imports: [ImageUpload],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  protected title = 'frontend';



}
