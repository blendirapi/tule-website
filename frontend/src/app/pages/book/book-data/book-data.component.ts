import { Component, Output, EventEmitter, Input, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { catchError } from 'rxjs/operators';
import { of } from 'rxjs';

interface BookingFormData {
  artist: string | null;
  service: string | null;
  date: string | null;
  time: string | null;
  firstName: string;
  lastName: string;
  phone: string;
  email: string;
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
export class BookDataComponent implements OnInit {
  @Input() selectedArtist: string | null = null;
  @Input() selectedService: string | null = null;
  @Input() selectedDate: Date | null = null;
  @Input() selectedTime: string | null = null;

  @Output() formSubmitted = new EventEmitter<boolean>();
  @Output() back = new EventEmitter<void>();

  firstName: string = '';
  lastName: string = '';
  phone: string = '';
  email: string = '';
  serviceName: string = '';
  artistName: string = '';

  constructor(private http: HttpClient) { }

  ngOnInit(): void {
    this.fetchNames();
  }

  private formatLocalDate(date: Date | null): string {
    if (!date) return '';

    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');

    return `${year}-${month}-${day}`;
  }

  fetchNames(): void {
    if (!this.selectedArtist || !this.selectedService) {
      return;
    }

    const params = {
      artist_id: this.selectedArtist,
      service_id: this.selectedService
    };

    this.http.get<NamesResponse>('/v0/api/service_artist_names', { params })
      .pipe(
        catchError((error: HttpErrorResponse) => {
          console.error('API error:', error);
          return of(null);
        })
      )
      .subscribe({
        next: (response) => {
          if (response) {
            this.serviceName = response.service_name || '';
            this.artistName = response.user_name || '';
          }
        }
      });
  }

  onSubmit(): void {
    const formData: BookingFormData = {
      artist: this.selectedArtist,
      service: this.selectedService,
      date: this.formatLocalDate(this.selectedDate),
      time: this.selectedTime,
      firstName: this.firstName,
      lastName: this.lastName,
      phone: this.phone,
      email: this.email
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

  goBack(): void {
    this.back.emit();
}   
}