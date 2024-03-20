import { goto as gotoNavigate } from "$app/navigation";
import { base } from '$app/paths';

// The 'goto' function is used to navigate to a specified
// destination within the application. It takes in two parameters:
// - destination: A string representing the path to navigate to.
// - opt?: object: An optional object that can contain additional navigation options.
export function goto(destination: string, opt?: object ) {
    const url = `${base}${destination}`;
    void gotoNavigate(url, opt);
}