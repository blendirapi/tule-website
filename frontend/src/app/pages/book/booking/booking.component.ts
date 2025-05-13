import { Component } from '@angular/core';
import { BookArtistComponent } from '../book-artist/book-artist.component';
import { BookTypeComponent } from '../book-type/book-type.component';
import { BookCalendarComponent } from '../book-calendar/book-calendar.component';
import { BookTimeComponent } from '../book-time/book-time.component';
import { BookDataComponent } from '../book-data/book-data.component';
import { BookFinalComponent } from '../book-final/book-final.component';
import { HeaderComponent } from '../../../components/header/header.component';
import { FooterComponent } from '../../../components/footer/footer.component';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-booking',
  standalone: true,
  imports: [
    HeaderComponent,
    FooterComponent,
    CommonModule,
    BookArtistComponent,
    BookTypeComponent,
    BookCalendarComponent,
    BookTimeComponent,
    BookDataComponent,
    BookFinalComponent
  ],
  templateUrl: './booking.component.html',
})
export class BookingComponent {
  selectedArtist: string | null = null;
  selectedArtistName: string | null = null;
  selectedService: string | null = null;
  selectedServiceName: string | null = null;
  price: number | null = null;
  bath: string | null = '';
  selectedDate: Date | null = null;
  selectedTime: string | null = null;
  bookingComplete = false;

  onArtistSelected(artist: string): void {
    this.selectedArtist = artist;
  }

  onArtistSelectedName(name: string): void {
    this.selectedArtistName = name;
  }

  onServiceSelected(service: string): void {
    this.selectedService = service;
  }

  onServiceSelectedName(name: string): void {
    this.selectedServiceName = name;
  }

  onServicePrice(price: number): void {
    this.price = price;
  }

  onServiceHasBath(hasBath: string): void {
    this.bath = hasBath;
  }

  onDateSelected(date: Date): void {
    this.selectedDate = date;
  }

  onTimeSelected(time: string): void {
    this.selectedTime = time;
  }

  onFormSubmitted(success: boolean): void {
    if (success) {
      this.bookingComplete = true;
    }
  }

  goBack(): void {
    if (this.selectedTime) {
      this.selectedTime = null;
    } else if (this.selectedDate) {
      this.selectedDate = null;
    } else if (this.selectedService) {
      this.selectedService = null;
    } else if (this.selectedArtist) {
      this.selectedArtist = null;
    }
  } 
}