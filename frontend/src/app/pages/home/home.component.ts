import {
  AfterViewInit,
  Component,
  CUSTOM_ELEMENTS_SCHEMA,
} from '@angular/core';
import { BannerComponent } from '../../components/banner/banner.component';
import { BookButtonComponent } from '../../components/book-button/book-button.component';
import { HeaderComponent } from '../../components/header/header.component';
import { FooterComponent } from '../../components/footer/footer.component';
import { CommonModule, NgStyle } from '@angular/common';
import Swiper from 'swiper';
import { Navigation, Scrollbar, Autoplay } from 'swiper/modules';

Swiper.use([Navigation, Scrollbar, Autoplay]);

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [
    HeaderComponent,
    FooterComponent,
    BannerComponent,
    BookButtonComponent,
    NgStyle,
    CommonModule,
  ],
  templateUrl: './home.component.html',
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
})
export class HomeComponent implements AfterViewInit {
  // PARALLAX
  // translateY = 0;

  // @HostListener('window:scroll', [])
  // onScroll(): void {
  //   const scrollTop = window.scrollY;
  //   this.translateY = scrollTop * -0.05;
  // }

  slidesFirstSwiper = [
    { title: 'Slide 1', image: 'banner.webp' },
    { title: 'Slide 2', image: 'icon.png' },
    { title: 'Slide 3', image: 'banner.webp' },
    { title: 'Slide 4', image: 'icon.png' },
    { title: 'Slide 5', image: 'banner.webp' },
    { title: 'Slide 6', image: 'icon.png' },
  ];

  slidesSecondSwiper = [
    { title: 'Slide 1', image: 'dummy-1.jpg' },
    { title: 'Slide 2', image: 'dummy-2.jpg' },
    { title: 'Slide 3', image: 'dummy-1.jpg' },
    { title: 'Slide 4', image: 'dummy-2.jpg' },
    { title: 'Slide 5', image: 'dummy-1.jpg' },
    { title: 'Slide 6', image: 'dummy-2.jpg' },
  ];

  ngAfterViewInit(): void {
    new Swiper('.firstSwiper', {
      slidesPerView: 1,
      loop: true,
      centeredSlides: false,
      autoplay: {
        delay: 3000,
        disableOnInteraction: false,
        reverseDirection: true,
      },
      scrollbar: {
        el: '.swiper-scrollbar',
        draggable: true,
      },
      navigation: {
        nextEl: '.swiper-button-next',
        prevEl: '.swiper-button-prev',
      },
    });

    new Swiper('.secondSwiper', {
      slidesPerView: 1,
      loop: true,
      centeredSlides: false,
      autoplay: {
        delay: 3000,
        disableOnInteraction: false,
      },
      scrollbar: {
        el: '.swiper-scrollbar.second',
        draggable: true,
      },
      navigation: {
        nextEl: '.swiper-button-next.second',
        prevEl: '.swiper-button-prev.second',
      },
    });
  }
}
