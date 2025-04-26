import { Component } from '@angular/core';
import { BannerComponent } from '../../components/banner/banner.component';
import { BookButtonComponent } from '../../components/book-button/book-button.component';
import { HeaderComponent } from '../../components/header/header.component';
import { FooterComponent } from '../../components/footer/footer.component';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [HeaderComponent, FooterComponent, BannerComponent, BookButtonComponent],
  templateUrl: './home.component.html',
})
export class HomeComponent {

}