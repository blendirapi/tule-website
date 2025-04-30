import { Component, HostListener } from '@angular/core';
import { BannerComponent } from '../../components/banner/banner.component';
import { BookButtonComponent } from '../../components/book-button/book-button.component';
import { HeaderComponent } from '../../components/header/header.component';
import { FooterComponent } from '../../components/footer/footer.component';
import { NgStyle } from '@angular/common';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [
    HeaderComponent,
    FooterComponent,
    BannerComponent,
    BookButtonComponent,
    NgStyle,
  ],
  templateUrl: './home.component.html',
})
export class HomeComponent {
  // PARALLAX
  translateY = 0;

  @HostListener('window:scroll', [])
  onScroll(): void {
    const scrollTop = window.scrollY;
    this.translateY = scrollTop * -0.05;
  }
}
