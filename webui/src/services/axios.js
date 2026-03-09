import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 30
});

// Redirect to login on 401 responses
instance.interceptors.response.use(
	response => response,
	error => {
		if (error.response && error.response.status === 401) {
			localStorage.removeItem("token");
			localStorage.removeItem("name");
			window.location.hash = "#/";
		}
		return Promise.reject(error);
	}
);

export default instance;
