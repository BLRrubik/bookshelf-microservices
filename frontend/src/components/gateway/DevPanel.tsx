import { useState } from 'react';
import { Terminal, Copy, Check, Trash2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Label } from '@/components/ui/label';
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet';
import { useRequestStore } from '@/stores/request';
import { useRateLimitStore } from '@/stores/rateLimit';

export function DevPanel() {
  const [open, setOpen] = useState(false);
  const [copied, setCopied] = useState(false);
  const { lastRequest, history, clearHistory } = useRequestStore();
  const rateLimit = useRateLimitStore();

  const handleCopyRequestId = () => {
    if (lastRequest?.requestId) {
      navigator.clipboard.writeText(lastRequest.requestId);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>
        <Button 
          variant="outline" 
          size="icon" 
          className="fixed bottom-4 right-4 z-40 shadow-lg"
        >
          <Terminal className="h-4 w-4" />
        </Button>
      </SheetTrigger>
      <SheetContent className="overflow-y-auto">
        <SheetHeader>
          <SheetTitle className="flex items-center gap-2">
            <Terminal className="h-5 w-5" />
            Developer Panel
          </SheetTitle>
        </SheetHeader>
        
        <div className="space-y-6 mt-6">
          {/* Last Request */}
          <div className="space-y-3">
            <Label className="text-muted-foreground uppercase text-xs tracking-wide">
              Последний запрос
            </Label>
            
            {lastRequest ? (
              <div className="space-y-3 p-3 bg-muted rounded-lg">
                <div className="flex items-center justify-between">
                  <Badge variant="outline">
                    {lastRequest.method}
                  </Badge>
                  <Badge variant={lastRequest.status < 400 ? 'success' : 'destructive'}>
                    {lastRequest.status}
                  </Badge>
                </div>
                
                <code className="block text-xs bg-background p-2 rounded break-all">
                  {lastRequest.url}
                </code>
                
                <div className="grid grid-cols-2 gap-2 text-sm">
                  <div>
                    <Label className="text-xs text-muted-foreground">X-Request-ID</Label>
                    <div className="flex items-center gap-1">
                      <code className="text-xs truncate flex-1">
                        {lastRequest.requestId || 'N/A'}
                      </code>
                      {lastRequest.requestId && (
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-6 w-6"
                          onClick={handleCopyRequestId}
                        >
                          {copied ? (
                            <Check className="h-3 w-3 text-green-500" />
                          ) : (
                            <Copy className="h-3 w-3" />
                          )}
                        </Button>
                      )}
                    </div>
                  </div>
                  
                  <div>
                    <Label className="text-xs text-muted-foreground">X-Cache</Label>
                    <Badge 
                      variant={lastRequest.cacheStatus === 'HIT' ? 'success' : 'outline'}
                      className="mt-1"
                    >
                      {lastRequest.cacheStatus || 'N/A'}
                    </Badge>
                  </div>
                </div>
                
                <div>
                  <Label className="text-xs text-muted-foreground">Response Time</Label>
                  <p className="text-sm font-mono">{lastRequest.responseTime}ms</p>
                </div>
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">Нет данных</p>
            )}
          </div>

          <Separator />

          {/* Rate Limit */}
          <div className="space-y-3">
            <Label className="text-muted-foreground uppercase text-xs tracking-wide">
              Rate Limit
            </Label>
            
            <div className="p-3 bg-muted rounded-lg space-y-2">
              <div className="flex justify-between items-center">
                <span className="text-sm">Remaining</span>
                <Badge variant={rateLimit.remaining > 10 ? 'outline' : 'destructive'}>
                  {rateLimit.remaining} / {rateLimit.limit}
                </Badge>
              </div>
              
              {rateLimit.reset && (
                <div className="flex justify-between items-center text-sm">
                  <span className="text-muted-foreground">Reset</span>
                  <span className="font-mono">
                    {new Date(rateLimit.reset * 1000).toLocaleTimeString('ru-RU')}
                  </span>
                </div>
              )}
              
              {rateLimit.isBlocked && (
                <Badge variant="destructive" className="w-full justify-center">
                  Blocked - Retry in {rateLimit.retryAfter}s
                </Badge>
              )}
            </div>
          </div>

          <Separator />

          {/* Request History */}
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <Label className="text-muted-foreground uppercase text-xs tracking-wide">
                История ({history.length})
              </Label>
              {history.length > 0 && (
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={clearHistory}
                >
                  <Trash2 className="h-3 w-3 mr-1" />
                  Очистить
                </Button>
              )}
            </div>
            
            <div className="space-y-2 max-h-64 overflow-y-auto">
              {history.slice(0, 10).map((req, i) => (
                <div 
                  key={i} 
                  className="flex items-center gap-2 text-xs p-2 bg-muted/50 rounded"
                >
                  <Badge variant="outline" className="text-[10px] px-1">
                    {req.method}
                  </Badge>
                  <span className="flex-1 truncate font-mono">{req.url}</span>
                  <Badge 
                    variant={req.cacheStatus === 'HIT' ? 'success' : 'secondary'}
                    className="text-[10px] px-1"
                  >
                    {req.cacheStatus || '-'}
                  </Badge>
                  <span className="text-muted-foreground w-12 text-right">
                    {req.responseTime}ms
                  </span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}





