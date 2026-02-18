'use client'
import SignUpForm from "../components/Signup"
function Signup() {
    const route = "http://localhost:1323/signup"
    console.log(route)
    return (

        <SignUpForm routes={route} destination="/login" />
    );
}
export default Signup   