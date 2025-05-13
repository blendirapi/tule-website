import { Routes } from '@angular/router';
import { HomeComponent } from './pages/home/home.component';
import { BookingComponent } from './pages/book/booking/booking.component';
import { LoginComponent } from './pages/login/login.component';
import { DailyCalendarComponent } from './pages/view-bookings/view-bookings.component';
import { UsersComponent } from './pages/users/users.component';
import { ServicesComponent } from './pages/services/services.component';
import { PrivacyPolicyComponent } from './pages/privacy-policy/privacy-policy.component';
import { TermsAndConditionsComponent } from './pages/terms-and-conditions/terms-and-conditions.component';
import { AuthGuard } from './auth/auth.guard';

export const routes: Routes = [
  { path: '', component: HomeComponent, title: 'Άρχικη | Tule Hood' },
  { path: 'booking', component: BookingComponent, title: 'Κλείσε ραντεβού | Tule Hood' },
  { path: 'login', component: LoginComponent, title: 'Σύνδεση | Tule Hood' },
  { path: 'dashboard', component: DailyCalendarComponent, canActivate: [AuthGuard], title: 'Κρατήσεις | Tule Hood' },
  { path: 'dashboard/users', component: UsersComponent, canActivate: [AuthGuard], title: 'Χρήστες | Tule Hood' },
  { path: 'dashboard/services', component: ServicesComponent, canActivate: [AuthGuard], title: 'Υπηρεσίες | Tule Hood' },
  { path: 'privacy', component: PrivacyPolicyComponent, title: 'Πολιτική Απορρήτου | Tule Hood' },
  { path: 'terms', component: TermsAndConditionsComponent, title: 'Όροι και Προϋποθέσεις | Tule Hood' },
  { path: '**', component: HomeComponent, title: 'Homepage | Tule Hood' },
];
