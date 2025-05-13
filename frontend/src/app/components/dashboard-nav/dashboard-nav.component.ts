import { Component } from '@angular/core';
import { AuthService } from '../../auth/auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-dashboard-nav',
  standalone: true,
  imports: [],
  templateUrl: './dashboard-nav.component.html',
})
export class DashboardNavComponent {
  constructor(private authService: AuthService, private router: Router) {}

	logout() {
		this.authService.logout();
		this.router.navigate(['/login']);
	}
}