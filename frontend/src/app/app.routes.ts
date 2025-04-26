import { Routes } from '@angular/router';
import { HomeComponent } from './pages/home/home.component';
import { BookingComponent } from './pages/book/booking/booking.component';
import { LoginComponent } from './pages/login/login.component';
import { DailyCalendarComponent } from './pages/view-bookings/view-bookings.component';
import { DashboardComponent } from './pages/dashboard/dashboard.component';

export const routes: Routes = [
  { path: '', component: HomeComponent, title: 'Homepage | Tule Hood' },
  { path: 'booking', component: BookingComponent, title: 'Booking | Tule Hood' },
  { path: 'login', component: LoginComponent, title: 'Login | Tule Hood' },
  { path: 'dashboard', component: DailyCalendarComponent, title: 'Dashboard | Tule Hood' },
  { path: 'dashboard/users', component: DashboardComponent, title: 'Dashboard | Tule Hood' },
  { path: '**', component: HomeComponent, title: 'Homepage | Tule Hood' },
];
