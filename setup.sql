
CREATE TABLE public.niveis_acesso (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(50) NOT NULL
);

INSERT INTO public.niveis_acesso (nome) VALUES ('user'), ('admin');

CREATE TABLE public.usuarios (
    id SERIAL PRIMARY KEY,
    nome_usuario VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    senha VARCHAR(255) NOT NULL,
    nivel_acesso_id INTEGER NOT NULL DEFAULT 1 REFERENCES public.niveis_acesso(id)
);
