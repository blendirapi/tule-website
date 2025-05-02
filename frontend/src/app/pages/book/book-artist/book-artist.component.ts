import { Component, Output, EventEmitter } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { catchError } from 'rxjs/operators';
import { of } from 'rxjs';
import { CommonModule } from '@angular/common';

interface Artist {
  user_id: string;
  name: string;
}

@Component({
  selector: 'app-book-artist',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './book-artist.component.html',
})
export class BookArtistComponent {
  @Output() artistSelected = new EventEmitter<string>();
  @Output() artistSelectedName = new EventEmitter<string>();
  artists: Artist[] = [];
  isLoading: boolean = true;

  constructor(private http: HttpClient) { }

  ngOnInit(): void {
    this.fetchArtists();
  }

  selectArtist(artistId: string, artistName: string) {
    this.artistSelected.emit(artistId);
    this.artistSelectedName.emit(artistName);
  }

  fetchArtists(): void {
    this.http.get<Artist[]>('/v0/api/artists')
      .pipe(
        catchError((error: HttpErrorResponse) => {
          console.error('API error:', error);
          return of([]);
        })
      )
      .subscribe({
        next: (response) => {
          this.artists = Array.isArray(response) ? response : [];
          this.isLoading = false;
        }
      });
  }
}