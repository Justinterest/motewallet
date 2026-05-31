export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="auth-shell flex min-h-screen flex-col items-center justify-center px-4">
      <div className="w-full max-w-[420px]">
        <div className="mb-8 text-center">
          <h1 className="text-3xl font-bold text-primary">Motewallet</h1>
          <p className="mt-2 text-sm text-muted-foreground">商户平台</p>
        </div>
        {children}
      </div>
    </div>
  );
}
