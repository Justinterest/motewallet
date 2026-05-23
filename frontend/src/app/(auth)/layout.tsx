export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-slate-50 px-4">
      <div className="w-full max-w-[420px]">
        <div className="mb-8 text-center">
          <h1 className="text-3xl font-bold text-blue-800">Motewallet</h1>
          <p className="mt-2 text-sm text-slate-500">商户平台</p>
        </div>
        {children}
      </div>
    </div>
  );
}
