import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { motion } from 'framer-motion';
import { User, BookOpen, Star } from 'lucide-react';
import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Separator } from '@/components/ui/separator';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { useAuthStore } from '@/stores/auth';
import { useUpdateProfile } from '@/api/auth';
import { useUserReviews } from '@/api/reviews';
import { formatDate, formatDistanceToNow } from '@/lib/date';

const profileSchema = z.object({
  username: z.string()
    .min(3, 'Минимум 3 символа')
    .max(50, 'Максимум 50 символов')
    .regex(/^[a-zA-Z0-9_]+$/, 'Только буквы, цифры и подчёркивания'),
});

type ProfileFormData = z.infer<typeof profileSchema>;

export function ProfilePage() {
  const user = useAuthStore((state) => state.user);
  const updateProfile = useUpdateProfile();
  const { data: userReviews, isLoading: reviewsLoading } = useUserReviews(user?.id || '');
  
  const form = useForm<ProfileFormData>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      username: user?.username || '',
    },
  });

  const onSubmit = (data: ProfileFormData) => {
    updateProfile.mutate(data);
  };

  const initials = user?.username
    ?.split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2) || 'U';

  if (!user) {
    return null;
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="max-w-2xl mx-auto space-y-6"
    >
      <h1 className="font-display text-3xl font-bold">Профиль</h1>

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
              <CardDescription>{user.email}</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <User className="h-4 w-4" />
            <span>Зарегистрирован: {formatDate(user.created_at)}</span>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Редактировать профиль</CardTitle>
          <CardDescription>
            Обновите информацию своего профиля
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                control={form.control}
                name="username"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Имя пользователя</FormLabel>
                    <FormControl>
                      <Input {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <div className="flex items-center gap-4">
                <FormItem className="flex-1">
                  <FormLabel>Email</FormLabel>
                  <Input value={user.email} disabled />
                </FormItem>
              </div>

              <Separator />

              <Button 
                type="submit" 
                disabled={updateProfile.isPending || !form.formState.isDirty}
              >
                {updateProfile.isPending ? 'Сохранение...' : 'Сохранить изменения'}
              </Button>
            </form>
          </Form>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Star className="h-5 w-5" />
            Мои рецензии
          </CardTitle>
          <CardDescription>
            Рецензии, которые вы оставили на книги
          </CardDescription>
        </CardHeader>
        <CardContent>
          {reviewsLoading ? (
            <p className="text-sm text-muted-foreground">Загрузка...</p>
          ) : !userReviews?.data?.length ? (
            <p className="text-sm text-muted-foreground">
              Вы ещё не оставили ни одной рецензии
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

