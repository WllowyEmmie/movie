'use client'
import LoginForm from "../components/Login";
function Login() {
    const route = "http://localhost:1323/login";
    console.log(route)
    return (

        <LoginForm routes={route} destination="" />
    );
}
export default Login;