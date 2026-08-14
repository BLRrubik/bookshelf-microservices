import { Lock, BookOpen, ArrowRight, Zap } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './card';
import { Button } from './button';

interface FeatureLockedProps {
  title: string;
  description: string;
  stage: number;
  hint?: string;
  icon?: React.ReactNode;
}

export function FeatureLocked({ 
  title, 
  description, 
  stage, 
  hint,
  icon
}: FeatureLockedProps) {
  return (
    <Card className="border-dashed border-2 border-muted-foreground/25 bg-muted/5">
      <CardHeader className="text-center pb-2">
        <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-muted">
          {icon || <Lock className="h-8 w-8 text-muted-foreground" />}
        </div>
        <CardTitle className="text-xl">{title}</CardTitle>
        <CardDescription className="text-base">
          {description}
        </CardDescription>
      </CardHeader>
      <CardContent className="text-center space-y-4">
        {stage > 0 && (
          <div className="inline-flex items-center gap-2 bg-primary/10 text-primary px-4 py-2 rounded-full text-sm font-medium">
            <BookOpen className="h-4 w-4" />
            <span>Смотри Этап {stage} в DETAILED_STAGES.md</span>
          </div>
        )}
        
        {hint && (
          <p className="text-sm text-muted-foreground max-w-md mx-auto">
            💡 {hint}
          </p>
        )}

        <div className="pt-4">
          <Button variant="outline" disabled className="gap-2">
            <span>Скоро будет доступно</span>
            <ArrowRight className="h-4 w-4" />
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}

interface GatewayNotRunningProps {
  onRetry?: () => void;
}

export function GatewayNotRunning({ onRetry }: GatewayNotRunningProps) {
  return (
    <Card className="border-destructive/50 bg-destructive/5">
      <CardHeader className="text-center pb-2">
        <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-destructive/10">
          <Zap className="h-8 w-8 text-destructive" />
        </div>
        <CardTitle className="text-xl text-destructive">API Gateway недоступен</CardTitle>
        <CardDescription className="text-base">
          Не удалось подключиться к API Gateway. Убедись, что он запущен.
        </CardDescription>
      </CardHeader>
      <CardContent className="text-center space-y-4">
        <div className="bg-muted p-4 rounded-lg text-left max-w-md mx-auto space-y-2">
          <p className="text-xs text-muted-foreground mb-2">Запусти api-gateway:</p>
          <p className="text-sm font-mono text-muted-foreground">
            $ cd api-gateway && go run ./cmd/server
          </p>
          <p className="text-xs text-muted-foreground mt-2">
            Gateway должен быть доступен на порту 8000
          </p>
        </div>
        
        <div className="inline-flex items-center gap-2 bg-primary/10 text-primary px-4 py-2 rounded-full text-sm font-medium">
          <BookOpen className="h-4 w-4" />
          <span>Смотри Этап 4 — Проксирующие хендлеры</span>
        </div>
        
        {onRetry && (
          <div className="pt-2">
            <Button onClick={onRetry} variant="outline">
              Попробовать снова
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  );
}

interface FeatureErrorProps {
  title: string;
  error?: Error | null;
  onRetry?: () => void;
}

export function FeatureError({ title, error, onRetry }: FeatureErrorProps) {
  const isNetworkError = error?.message?.includes('Network Error') || 
                         error?.message?.includes('ERR_CONNECTION_REFUSED');
  
  if (isNetworkError) {
    return <GatewayNotRunning onRetry={onRetry} />;
  }
  
  return (
    <Card className="border-destructive/50 bg-destructive/5">
      <CardHeader className="text-center pb-2">
        <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-destructive/10">
          <span className="text-3xl">⚠️</span>
        </div>
        <CardTitle className="text-xl text-destructive">{title}</CardTitle>
        <CardDescription className="text-base">
          Произошла ошибка при загрузке данных.
        </CardDescription>
      </CardHeader>
      <CardContent className="text-center space-y-4">
        {onRetry && (
          <Button onClick={onRetry} variant="outline">
            Попробовать снова
          </Button>
        )}
      </CardContent>
    </Card>
  );
}
