import { Component, Output, EventEmitter, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { catchError } from 'rxjs/operators';
import { of } from 'rxjs';
import { registerLocaleData } from '@angular/common';
import localeEl from '@angular/common/locales/el';

registerLocaleData(localeEl);

interface BookingFormData {
  artist: string | null;
  service: string | null;
  date: string | null;
  time: string | null;
  firstName: string;
  lastName: string;
  phone: string;
  email: string;
  bath: string;
}

interface NamesResponse {
  service_name?: string;
  user_name?: string;
}

interface BookingResponse extends NamesResponse {
  success?: boolean;
  message?: string;
  booking_id: number;
}

@Component({
  selector: 'app-book-data',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './book-data.component.html',
})
export class BookDataComponent {
  @Input() selectedArtist: string | null = null;
  @Input() selectedArtistName: string | null = null;
  @Input() selectedService: string | null = null;
  @Input() selectedServiceName: string | null = null;
  @Input() selectedDate: Date | null = null;
  @Input() selectedTime: string | null = null;
  @Input() price: number | null = null;
  @Input() hasBath: string | '' = '';

  @Output() formSubmitted = new EventEmitter<boolean>();

  firstName: string = '';
  lastName: string = '';
  phone: string = '';
  email: string = '';
  serviceName: string = '';
  isSubmitting: boolean = false;
  honeypot: string = '';

  constructor(private http: HttpClient) { }

  private formatLocalDate(date: Date | null): string {
    if (!date) return '';

    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');

    return `${year}-${month}-${day}`;
  }

  onSubmit(): void {
    if (this.honeypot) {
      return;
    }

    this.isSubmitting = true;
    const formData: BookingFormData = {
      artist: this.selectedArtist,
      service: this.selectedService,
      date: this.formatLocalDate(this.selectedDate),
      time: this.selectedTime,
      firstName: this.firstName,
      lastName: this.lastName,
      phone: this.phone,
      email: this.email,
      bath: this.hasBath ? "true" : "false"
    };

    this.http.post<BookingResponse>('/v0/api/add_booking', formData)
      .pipe(
        catchError((error: HttpErrorResponse) => {
          console.error('API error:', error);
          return of(null);
        })
      )
      .subscribe({
        next: (response) => {
          if (response && response.booking_id > 0) {
            this.formSubmitted.emit(true);
          }
        }
      });
  }
}