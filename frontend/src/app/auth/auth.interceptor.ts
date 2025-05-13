import { HttpInterceptorFn, HttpRequest, HttpHandlerFn } from '@angular/common/http';

export const AuthInterceptor: HttpInterceptorFn = (req: HttpRequest<any>, next: HttpHandlerFn) => {
	const token = localStorage.getItem('auth');

	if (token) {
		const cloned = req.clone({
			headers: req.headers.set('Authorization', token)
		});
		return next(cloned);
	}

	return next(req);
};
