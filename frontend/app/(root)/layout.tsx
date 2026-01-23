import Header from "@/components/Header"
import { getCurrentUser } from "@/lib/actions/auth.actions"

const Layout = async ({ children }: {children: React.ReactNode}) => {
  const user = await getCurrentUser();
  
  return (
    <main className = "min-h-screen text-gray-400">
		  <Header user={user}/>
    <div className = "container py-10">
			  {children}
    </div>

    </main>
  )
}

export default Layout
