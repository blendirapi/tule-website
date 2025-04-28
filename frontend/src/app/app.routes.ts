import { Routes } from '@angular/router';
import { HomeComponent } from './pages/home/home.component';
import { BookingComponent } from './pages/book/booking/booking.component';
import { LoginComponent } from './pages/login/login.component';
import { DailyCalendarComponent } from './pages/view-bookings/view-bookings.component';
import { UsersComponent } from './pages/users/users.component';
import { ServicesComponent } from './pages/services/services.component';

export const routes: Routes = [
  { path: '', component: HomeComponent, title: 'Homepage | Tule Hood' },
  { path: 'booking', component: BookingComponent, title: 'Booking | Tule Hood' },
  { path: 'login', component: LoginComponent, title: 'Login | Tule Hood' },
  { path: 'dashboard', component: DailyCalendarComponent, title: 'Dashboard | Tule Hood' },
  { path: 'dashboard/users', component: UsersComponent, title: 'Dashboard | Tule Hood' },
  { path: 'dashboard/services', component: ServicesComponent, title: 'Dashboard | Tule Hood' },
  { path: '**', component: HomeComponent, title: 'Homepage | Tule Hood' },
];
