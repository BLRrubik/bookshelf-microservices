import { useParams, Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { User, BookOpen, Star, ArrowLeft } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { useUser } from '@/api/auth';
import { useUserReviews } from '@/api/reviews';
import { formatDate, formatDistanceToNow } from '@/lib/date';

export function UserProfilePage() {
  const { userId } = useParams<{ userId: string }>();
  const { data: user, isLoading: userLoading, error: userError } = useUser(userId || '');
  const { data: userReviews, isLoading: reviewsLoading } = useUserReviews(userId || '');

  if (userLoading) {
    return (
      <div className="max-w-2xl mx-auto">
        <p className="text-muted-foreground">Загрузка...</p>
      </div>
    );
  }

  if (userError || !user) {
    return (
      <div className="max-w-2xl mx-auto space-y-4">
        <Button variant="ghost" asChild>
          <Link to="/">
            <ArrowLeft className="h-4 w-4 mr-2" />
            На главную
          </Link>
        </Button>
        <Card>
          <CardContent className="py-8 text-center">
            <p className="text-muted-foreground">Пользователь не найден</p>
          </CardContent>
        </Card>
      </div>
    );
  }

  const initials = user.username
    .split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2);

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="max-w-2xl mx-auto space-y-6"
    >
      <Button variant="ghost" asChild>
        <Link to="/">
          <ArrowLeft className="h-4 w-4 mr-2" />
          На главную
        </Link>
      </Button>

      <Card>
        <CardHeader>
          <div className="flex items-center gap-4">
            <Avatar className="h-16 w-16">
              <AvatarFallback className="bg-primary/10 text-primary text-xl">
                {initials}
              </AvatarFallback>
            </Avatar>
            <div>
              <CardTitle>{user.username}</CardTitle>
              <CardDescription className="flex items-center gap-2 mt-1">
                <User className="h-4 w-4" />
                Зарегистрирован: {formatDate(user.created_at)}
              </CardDescription>
            </div>
          </div>
        </CardHeader>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Star className="h-5 w-5" />
            Рецензии пользователя
          </CardTitle>
          <CardDescription>
            Рецензии, оставленные {user.username}
          </CardDescription>
        </CardHeader>
        <CardContent>
          {reviewsLoading ? (
            <p className="text-sm text-muted-foreground">Загрузка...</p>
          ) : !userReviews?.data?.length ? (
            <p className="text-sm text-muted-foreground">
              Пользователь ещё не оставил ни одной рецензии
            </p>
          ) : (
            <div className="space-y-4">
              {userReviews.data.map((review: any) => (
                <div key={review.id} className="border-b pb-4 last:border-0 last:pb-0">
                  <div className="flex items-start justify-between gap-2 mb-2">
                    <Link 
                      to={`/books/${review.book?.id || review.book_id}`}
                      className="flex items-center gap-2 text-sm font-medium hover:text-primary transition-colors"
                    >
                      <BookOpen className="h-4 w-4" />
                      {review.book?.title || 'Книга'}
                    </Link>
                    <div className="flex items-center gap-1">
                      {[...Array(5)].map((_, i) => (
                        <Star
                          key={i}
                          className={`h-3 w-3 ${
                            i < review.rating
                              ? 'fill-yellow-400 text-yellow-400'
                              : 'text-muted-foreground'
                          }`}
                        />
                      ))}
                    </div>
                  </div>
                  {review.title && (
                    <h4 className="font-medium text-sm mb-1">{review.title}</h4>
                  )}
                  <p className="text-sm text-muted-foreground line-clamp-2">
                    {review.content}
                  </p>
                  <time className="text-xs text-muted-foreground mt-2 block">
                    {formatDistanceToNow(review.created_at)}
                  </time>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </motion.div>
  );
}

