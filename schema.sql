-- Create Enums
CREATE TYPE difficulty_level AS ENUM ('beginner', 'intermediate', 'advanced');
CREATE TYPE question_type AS ENUM ('choice', 'fill_blank');

-- Create Profiles table (extends auth.users)
CREATE TABLE profiles (
  id UUID REFERENCES auth.users(id) PRIMARY KEY,
  username TEXT,
  avatar_url TEXT,
  current_level difficulty_level DEFAULT 'beginner',
  created_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc'::text, now())
);

-- Create Courses table
CREATE TABLE courses (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  level difficulty_level NOT NULL,
  content_md TEXT,
  order_index INTEGER,
  slug TEXT UNIQUE NOT NULL
);

-- Create Exercises table
CREATE TABLE exercises (
  id SERIAL PRIMARY KEY,
  course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE,
  question TEXT NOT NULL,
  options JSONB, -- For 'choice' type
  correct_answer TEXT NOT NULL,
  explanation TEXT,
  type question_type NOT NULL
);

-- Create User Progress table
CREATE TABLE user_progress (
  id SERIAL PRIMARY KEY,
  user_id UUID REFERENCES profiles(id) ON DELETE CASCADE,
  course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE,
  is_completed BOOLEAN DEFAULT FALSE,
  score INTEGER,
  completed_at TIMESTAMP WITH TIME ZONE DEFAULT timezone('utc'::text, now())
);

-- Row Level Security (RLS) Policies
-- Enable RLS
ALTER TABLE profiles ENABLE ROW LEVEL SECURITY;
ALTER TABLE courses ENABLE ROW LEVEL SECURITY;
ALTER TABLE exercises ENABLE ROW LEVEL SECURITY;
ALTER TABLE user_progress ENABLE ROW LEVEL SECURITY;

-- Profiles: Users can read all profiles, but only update their own.
CREATE POLICY "Public profiles are viewable by everyone." ON profiles FOR SELECT USING (true);
CREATE POLICY "Users can update own profile." ON profiles FOR UPDATE USING (auth.uid() = id);
-- Note: Profiles are usually created via a trigger on auth.users, so insert policy might be needed if created manually,
-- but typically we rely on triggers. For simplicity, we allow insert if it matches auth uid.
CREATE POLICY "Users can insert own profile." ON profiles FOR INSERT WITH CHECK (auth.uid() = id);

-- Courses: Publicly readable. Only admins (service role) can insert/update/delete.
CREATE POLICY "Courses are viewable by everyone." ON courses FOR SELECT USING (true);

-- Exercises: Publicly readable.
CREATE POLICY "Exercises are viewable by everyone." ON exercises FOR SELECT USING (true);

-- User Progress: Users can view and update their own progress.
CREATE POLICY "Users can view own progress." ON user_progress FOR SELECT USING (auth.uid() = user_id);
CREATE POLICY "Users can insert own progress." ON user_progress FOR INSERT WITH CHECK (auth.uid() = user_id);
CREATE POLICY "Users can update own progress." ON user_progress FOR UPDATE USING (auth.uid() = user_id);

-- Function and Trigger to automatically create profile on signup
CREATE OR REPLACE FUNCTION public.handle_new_user()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO public.profiles (id, username, avatar_url)
  VALUES (new.id, new.raw_user_meta_data->>'username', new.raw_user_meta_data->>'avatar_url');
  RETURN new;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE TRIGGER on_auth_user_created
  AFTER INSERT ON auth.users
  FOR EACH ROW EXECUTE PROCEDURE public.handle_new_user();
