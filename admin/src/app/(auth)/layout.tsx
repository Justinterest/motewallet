export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-100">
      <div className="w-full max-w-[400px] px-4">
        <div className="mb-8 text-center">
          <h1 className="text-2xl font-bold text-slate-900">
            Motewallet 管理后台
          </h1>
        </div>
        {children}
      </div>
    </div>
  );
}
