import { Component, Input, Output, EventEmitter, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpErrorResponse, HttpParams } from '@angular/common/http';
import { catchError, of } from 'rxjs';

interface AvailableTimesResponse {
    success: boolean;
    message: string;
    times: string[];
}

@Component({
    selector: 'app-book-time',
    standalone: true,
    imports: [CommonModule],
    templateUrl: './book-time.component.html',
})
export class BookTimeComponent implements OnInit {
    @Input() selectedDate!: Date;
    @Input() selectedArtist: string | null = null;
    @Output() timeSelected = new EventEmitter<string>();

    timeSlots: string[] = [];
    isLoading: boolean = true;

    constructor(private http: HttpClient) { }

    ngOnInit(): void {
        this.fetchAvailableTimes();
    }

    fetchAvailableTimes(): void {
        if (!this.selectedArtist) return;

        let params = new HttpParams()
            .set('date', this.formatLocalDate(this.selectedDate))
            .set('artist_id', this.selectedArtist);

        this.http.get<AvailableTimesResponse>('/v0/api/available_times', { params })
            .pipe(
                catchError((error: HttpErrorResponse) => {
                    console.error('GET API error:', error);
                    return of({ success: false, message: error.message, times: [] });
                })
            )
            .subscribe((response) => {
                if (response.success) {
                    this.timeSlots = response.times || [];
                    this.isLoading = false;
                }
            });
    }

    selectTime(time: string): void {
        this.timeSelected.emit(time);
    }

    private formatLocalDate(date: Date): string {
        if (!date) return '';

        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');

        return `${year}-${month}-${day}`;
    }
}