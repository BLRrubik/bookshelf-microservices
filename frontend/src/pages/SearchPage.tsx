import { useState, useEffect } from 'react';
import { useSearchParams, Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { Search, Star } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { useSearch } from '@/api/dashboard';

export function SearchPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const query = searchParams.get('q') || '';
  const [inputValue, setInputValue] = useState(query);

  const { data, isLoading, isFetching } = useSearch({ q: query });

  useEffect(() => {
    setInputValue(query);
  }, [query]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = inputValue.trim();
    if (trimmed) {
      setSearchParams({ q: trimmed });
    }
  };

  // Debounced search
  useEffect(() => {
    const timeout = setTimeout(() => {
      const trimmed = inputValue.trim();
      if (trimmed.length >= 2 && trimmed !== query) {
        setSearchParams({ q: trimmed });
      }
    }, 300);
    
    return () => clearTimeout(timeout);
  }, [inputValue]);

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="space-y-6"
    >
      <div className="max-w-2xl mx-auto">
        <h1 className="font-display text-3xl font-bold text-center mb-6">
          Поиск книг
        </h1>
        
        <form onSubmit={handleSubmit} className="relative">
          <Search className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-muted-foreground" />
          <Input
            type="search"
            placeholder="Введите название или автора..."
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            className="pl-12 h-12 text-lg"
            autoFocus
          />
        </form>
      </div>

      {/* Results */}
      {query && (
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <p className="text-muted-foreground">
              {isLoading ? 'Поиск...' : `Найдено ${data?.total || 0} результатов`}
            </p>
            {isFetching && !isLoading && (
              <Badge variant="outline">Обновление...</Badge>
            )}
          </div>

          {isLoading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {Array.from({ length: 4 }).map((_, i) => (
                <Skeleton key={i} className="h-24" />
              ))}
            </div>
          ) : data?.books.length === 0 ? (
            <Card>
              <CardContent className="py-12 text-center">
                <p className="text-muted-foreground">
                  По запросу "{query}" ничего не найдено
                </p>
              </CardContent>
            </Card>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {data?.books.map((book, index) => (
                <motion.div
                  key={book.id}
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: index * 0.05 }}
                >
                  <Link to={`/books/${book.id}`}>
                    <Card className="hover:border-primary/50 transition-colors">
                      <CardContent className="p-4 flex items-center gap-4">
                        {book.cover_url ? (
                          <img 
                            src={book.cover_url} 
                            alt={book.title}
                            className="w-16 h-20 object-cover rounded"
                          />
                        ) : (
                          <div className="w-16 h-20 bg-muted rounded flex items-center justify-center">
                            <Search className="h-6 w-6 text-muted-foreground" />
                          </div>
                        )}
                        <div className="flex-1 min-w-0">
                          <h3 className="font-medium truncate">{book.title}</h3>
                          <p className="text-sm text-muted-foreground truncate">
                            {book.author}
                          </p>
                          <div className="flex items-center gap-2 mt-2">
                            {book.average_rating && (
                              <Badge variant="outline" className="gap-1">
                                <Star className="h-3 w-3 fill-primary text-primary" />
                                {book.average_rating.toFixed(1)}
                              </Badge>
                            )}
                            <span className="text-xs text-muted-foreground">
                              {book.reviews_count} рецензий
                            </span>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  </Link>
                </motion.div>
              ))}
            </div>
          )}
        </div>
      )}

      {!query && (
        <Card>
          <CardContent className="py-12 text-center">
            <Search className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
            <p className="text-muted-foreground">
              Введите запрос для поиска книг
            </p>
          </CardContent>
        </Card>
      )}
    </motion.div>
  );
}





