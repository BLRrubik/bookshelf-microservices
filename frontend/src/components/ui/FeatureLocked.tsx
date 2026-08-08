import { Lock, BookOpen, ArrowRight, Server } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './card.tsx';
import { Button } from './button.tsx';

interface FeatureLockedProps {
  title: string;
  description: string;
  stage: number;
  hint?: string;
  icon?: React.ReactNode;
  serviceName?: string;
}

export function FeatureLocked({ 
  title, 
  description, 
  stage, 
  hint,
  icon,
  serviceName 
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
        {serviceName && (
          <div className="inline-flex items-center gap-2 bg-blue-500/10 text-blue-400 px-4 py-2 rounded-full text-sm font-medium">
            <Server className="h-4 w-4" />
            <span>{serviceName}</span>
          </div>
        )}
        
        <div className="inline-flex items-center gap-2 bg-primary/10 text-primary px-4 py-2 rounded-full text-sm font-medium">
          <BookOpen className="h-4 w-4" />
          <span>Смотри Этап {stage} в DETAILED_STAGES.md</span>
        </div>
        
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

interface FeatureErrorProps {
  title: string;
  error?: Error | null;
  onRetry?: () => void;
  serviceName?: string;
  servicePort?: number;
}

export function FeatureError({ title, error, onRetry, serviceName, servicePort }: FeatureErrorProps) {
  const isNetworkError = error?.message?.includes('Network Error') || 
                         error?.message?.includes('ERR_CONNECTION_REFUSED');
  
  return (
    <Card className="border-destructive/50 bg-destructive/5">
      <CardHeader className="text-center pb-2">
        <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-full bg-destructive/10">
          <span className="text-3xl">⚠️</span>
        </div>
        <CardTitle className="text-xl text-destructive">{title}</CardTitle>
        <CardDescription className="text-base">
          {isNetworkError 
            ? `Не удалось подключиться к ${serviceName || 'серверу'}. Убедись, что сервис запущен.`
            : 'Произошла ошибка при загрузке данных.'
          }
        </CardDescription>
      </CardHeader>
      <CardContent className="text-center space-y-4">
        {isNetworkError && serviceName && (
          <div className="bg-muted p-4 rounded-lg text-left max-w-md mx-auto space-y-2">
            <p className="text-xs text-muted-foreground mb-2">Запусти {serviceName}:</p>
            <p className="text-sm font-mono text-muted-foreground">
              $ cd {serviceName} && go run ./cmd/server
            </p>
            {servicePort && (
              <p className="text-xs text-muted-foreground">
                Сервис должен быть доступен на порту {servicePort}
              </p>
            )}
          </div>
        )}
        
        {onRetry && (
          <Button onClick={onRetry} variant="outline">
            Попробовать снова
          </Button>
        )}
      </CardContent>
    </Card>
  );
}

interface ServiceStatusProps {
  services: Array<{
    name: string;
    status: 'online' | 'offline' | 'loading';
    port: number;
  }>;
}

export function ServiceStatusBar({ services }: ServiceStatusProps) {
  return (
    <div className="flex gap-4 justify-center py-2">
      {services.map((service) => (
        <div 
          key={service.name}
          className="flex items-center gap-2 text-sm"
        >
          <div 
            className={`w-2 h-2 rounded-full ${
              service.status === 'online' 
                ? 'bg-green-500' 
                : service.status === 'loading'
                ? 'bg-yellow-500 animate-pulse'
                : 'bg-red-500'
            }`}
          />
          <span className="text-muted-foreground">
            {service.name}
            <span className="text-xs ml-1">:{service.port}</span>
          </span>
        </div>
      ))}
    </div>
  );
}
