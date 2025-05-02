import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { DashboardNavComponent } from '../../components/dashboard-nav/dashboard-nav.component';

@Component({
  selector: 'app-services',
  standalone: true,
  imports: [CommonModule, FormsModule, DashboardNavComponent],
  templateUrl: './services.component.html',
})
export class ServicesComponent implements OnInit {
  services: any[] = [];
  form = { service_id: null, name: '', color: '', price: 0, time: '' };
  isEditing: Boolean = false;
  isSubmitting: Boolean = false;

  constructor(private http: HttpClient) { }

  ngOnInit() {
    this.loadServices();
  }

  loadServices() {
    this.http.get<any[]>('/v0/api/services').subscribe((data) => {
      this.services = data;
      this.isSubmitting = false;
    });
  }

  saveService() {
    this.isSubmitting = true;
    const apiUrl = this.isEditing ? `/v0/api/service/${this.form.service_id}` : '/v0/api/services';
    const method = this.isEditing ? 'put' : 'post';

    if(this.form.color == '') {
      this.form.color = '#000000';
    }

    this.http[method](apiUrl, this.form).subscribe(() => {
      this.loadServices();
      this.cancelEdit();
    });
  }

  editService(service: any) {
    this.form = { ...service };
    this.isEditing = true;
  }

  deleteService(serviceId: number) {
    if (confirm('Είστε σίγουροι ότι θέλετε να διαγράψετε αυτήν την υπηρεσία;')) {
      this.http.delete(`/v0/api/service/${serviceId}`).subscribe(() => this.loadServices());
    }
  }

  cancelEdit() {
    this.form = { service_id: null, name: '', color: '', price: 0, time: '' };
    this.isEditing = false;
  }
}