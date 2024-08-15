// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package html

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Login() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8\"><div class=\"max-w-md w-full\"><div><img class=\"mx-auto h-12 w-auto\" src=\"https://tailwindui.com/img/logos/v1/workflow-mark-on-white.svg\" alt=\"Workflow\"><h2 class=\"mt-6 text-center text-3xl leading-9 font-extrabold text-gray-900\">Sign in to your account</h2><p class=\"mt-2 text-center text-sm leading-5 text-gray-600\">We'll create one for you if it doesn't exist yet, don't you fucking worry about it.</p></div><div class=\"mt-6\"><button type=\"submit\" class=\"group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm leading-5 font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-500 focus:outline-none focus:border-indigo-700 focus:shadow-outline-indigo active:bg-indigo-700 transition duration-150 ease-in-out\" hx-get=\"/auth/google/login\" hx-target=\"body\" hx-push-url=\"true\"><span class=\"absolute left-0 inset-y-0 flex items-center pl-3\"><svg class=\"h-5 w-5 text-indigo-500 group-hover:text-indigo-400 transition ease-in-out duration-150\" fill=\"currentColor\" viewBox=\"0 0 20 20\"><path fill-rule=\"evenodd\" d=\"M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z\" clip-rule=\"evenodd\"></path></svg></span> Sign in with Google</button></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
